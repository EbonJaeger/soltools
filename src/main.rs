mod errors;

use crate::errors::{SolError, SolResult};

use fs_extra::copy_items;
use fs_extra::dir::CopyOptions;
use fs_extra::remove_items;
use git2::Repository;
use glob::glob;
use std::process;
use std::process::Command;
use std::{env, io::Write};
use std::{
    error::Error,
    path::{Path, PathBuf},
};
use std::{fs, fs::File};

const LOCAL_REPO_PATH: &str = "/var/lib/solbuild/local";

/// Copies any eopkg files in the current directory to the local solbuild
/// repo. This does not index the repo afterwards.
fn copy_packages() -> SolResult<()> {
    let mut packages = Vec::new();

    println!("Looking for packages to copy...");
    for entry in glob("*.eopkg").unwrap() {
        match entry {
            Ok(path) => {
                println!(
                    "Found package: {}",
                    path.file_name().unwrap().to_str().unwrap()
                );
                packages.push(path);
            }
            Err(e) => return Err(SolError::Glob(e)),
        }
    }

    let mut options = CopyOptions::new();
    options.overwrite = true;

    match copy_items(&packages, LOCAL_REPO_PATH, &options) {
        Ok(_) => Ok(()),
        Err(e) => Err(SolError::Fs(e)),
    }
}

/// Removes all eopkg files from the local solbuild repo. This
/// does not index the local repo afterwards.
fn clean_local_repo() -> SolResult<()> {
    let mut paths = Vec::new();
    let search = format!("{}/*.eopkg", LOCAL_REPO_PATH);

    for entry in glob(&search).unwrap() {
        match entry {
            Ok(path) => paths.push(path),
            Err(e) => return Err(SolError::Glob(e)),
        }
    }

    match remove_items(&paths) {
        Ok(_) => Ok(()),
        Err(e) => Err(SolError::Fs(e)),
    }
}

fn clone_repo<P: AsRef<Path>>(current_dir: P, package: &str) -> SolResult<()> {
    // Check that we're in the root packaging directory where `common` lives
    let mut common_path = PathBuf::new();
    common_path.push(&current_dir);
    common_path.push("common");

    if fs::metadata(&common_path).is_err() {
        return Err(SolError::Other(
            "not in packaging root directory: 'common' not found",
        ));
    }

    // Clone the repo
    let url = format!("https://dev.getsol.us/source/{}.git", package);
    match Repository::clone(&url, &package) {
        Ok(_) => Ok(()),
        Err(e) => Err(SolError::Git(e)),
    }
}

/// Indexes the packages in the local solbuild repository.
fn index_repo() -> SolResult<()> {
    let status = Command::new("eopkg")
        .current_dir(LOCAL_REPO_PATH)
        .arg("index")
        .arg("--skip-signing")
        .arg(LOCAL_REPO_PATH)
        .status();

    match status {
        Ok(_) => Ok(()),
        Err(e) => Err(SolError::Io(e)),
    }
}

fn init_repo<P: AsRef<Path>>(current_dir: P, name: &str, source_url: &str) -> SolResult<()> {
    // Check that we're in the root packaging directory where `common` lives
    let mut common_path = PathBuf::new();
    common_path.push(&current_dir);
    common_path.push("common");

    if fs::metadata(&common_path).is_err() {
        return Err(SolError::Other(
            "not in packaging root directory: 'common' not found",
        ));
    }

    // Create a new package directory
    let mut repo_path = PathBuf::new();
    repo_path.push(&current_dir);
    repo_path.push(name);

    if let Err(e) = fs::create_dir(&repo_path) {
        return Err(SolError::Io(e));
    }

    // Create the repo's Makefile
    println!("Creating package Makefile");
    let mut makefile_path = PathBuf::new();
    makefile_path.push(&repo_path);
    makefile_path.push("Makefile");

    let mut makefile = File::create(makefile_path)?;
    makefile.write_all(b"include ../Makefile.common\n")?;

    println!("\nRunning yauto.py");

    // Run yauto.py to generate the package.yml
    let yauto = format!("{}/Scripts/yauto.py", common_path.to_str().unwrap());
    Command::new(&yauto)
        .current_dir(&repo_path)
        .arg(source_url)
        .status()?;

    println!("\nCreating git repo");

    // Create a git repo in the new directory
    match Repository::init(&repo_path) {
        Ok(_) => Ok(()),
        Err(e) => Err(SolError::Git(e)),
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        eprintln!("Invalid arguments");
        process::exit(1);
    }

    let command = args[1].as_ref();
    match command {
        "copy" => {
            sudo::escalate_if_needed()?;

            if let Err(e) = copy_packages() {
                eprintln!("Error copying packages to local repo: {}", e);
                process::exit(1);
            }

            println!();

            if let Err(e) = index_repo() {
                eprintln!("Error indexing local repo: {}", e);
                process::exit(1);
            }

            Ok(())
        }
        "clean" => {
            sudo::escalate_if_needed()?;

            if let Err(e) = clean_local_repo() {
                eprintln!("Error cleaning local repo: {}", e);
                process::exit(1);
            }

            if let Err(e) = index_repo() {
                eprintln!("Error indexing local repo: {}", e);
                process::exit(1);
            }

            Ok(())
        }
        "clone" => {
            if args.len() != 3 {
                eprintln!("Invalid arguments. Usage: soltools clone NAME");
                process::exit(1);
            }

            let current_dir = match env::current_dir() {
                Ok(dir) => dir,
                Err(e) => {
                    eprintln!("Error getting current working directory: {}", e);
                    process::exit(1);
                }
            };

            if let Err(e) = clone_repo(current_dir, &args[2]) {
                eprintln!("Error creating new repo: {}", e);
                process::exit(1);
            }

            Ok(())
        }
        "init" => {
            if args.len() != 4 {
                eprintln!("Invalid arguments. Usage: soltools init NAME URL");
                process::exit(1);
            }

            let current_dir = match env::current_dir() {
                Ok(dir) => dir,
                Err(e) => {
                    eprintln!("Error getting current working directory: {}", e);
                    process::exit(1);
                }
            };

            if let Err(e) = init_repo(current_dir, &args[2], &args[3]) {
                eprintln!("Error creating new repo: {}", e);
                process::exit(1);
            }

            Ok(())
        }
        _ => {
            eprintln!("Unknown subcommand: {}", command);
            process::exit(1);
        }
    }
}
