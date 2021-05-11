use fs_extra::copy_items;
use fs_extra::dir::CopyOptions;
use fs_extra::remove_items;
use glob::glob;
use std::env;
use std::process;
use std::process::Command;
use std::{error::Error, path::PathBuf};

const LOCAL_REPO_PATH: &str = "/var/lib/solbuild/local";

/// Copies any eopkg files in the current directory to the local solbuild
/// repo. This does not index the repo afterwards.
fn copy_packages(current_dir: PathBuf) -> Result<(), Box<dyn Error>> {
    let mut packages = Vec::new();

    println!("Looking for packages to copy...");
    let search = format!("{}/*.eopkg", current_dir.to_str().unwrap_or("."));
    for entry in glob(&search).unwrap() {
        match entry {
            Ok(path) => {
                println!(
                    "Found package: {}",
                    path.file_name().unwrap().to_str().unwrap()
                );
                packages.push(path);
            }
            Err(e) => return Err(Box::new(e)),
        }
    }

    let mut options = CopyOptions::new();
    options.overwrite = true;

    match copy_items(&packages, LOCAL_REPO_PATH, &options) {
        Ok(_) => Ok(()),
        Err(e) => Err(Box::new(e)),
    }
}

/// Removes all eopkg files from the local solbuild repo. This
/// does not index the local repo afterwards.
fn clean_local_repo() -> Result<(), Box<dyn Error>> {
    let mut paths = Vec::new();
    let search = format!("{}/*.eopkg", LOCAL_REPO_PATH);

    for entry in glob(&search).unwrap() {
        match entry {
            Ok(path) => paths.push(path),
            Err(e) => return Err(Box::new(e)),
        }
    }

    match remove_items(&paths) {
        Ok(_) => Ok(()),
        Err(e) => Err(Box::new(e)),
    }
}

/// Indexes the packages in the local solbuild repository.
fn index_repo() -> Result<(), Box<dyn Error>> {
    let status = Command::new("eopkg")
        .current_dir(LOCAL_REPO_PATH)
        .arg("index")
        .arg("--skip-signing")
        .arg(LOCAL_REPO_PATH)
        .status();

    match status {
        Ok(_) => Ok(()),
        Err(e) => Err(Box::new(e)),
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() != 2 {
        eprintln!("Invalid arguments");
        process::exit(1);
    }

    let command = args[1].as_ref();
    match command {
        "copy" => {
            sudo::escalate_if_needed()?;

            let current_dir = match env::current_dir() {
                Ok(dir) => dir,
                Err(e) => {
                    eprintln!("Error getting current working directory: {}", e);
                    process::exit(1);
                }
            };

            if let Err(e) = copy_packages(current_dir) {
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
        _ => {
            eprintln!("Unknown subcommand: {}", command);
            process::exit(1);
        }
    }
}
