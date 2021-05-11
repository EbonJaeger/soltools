use fs_extra::copy_items;
use fs_extra::dir::CopyOptions;
use glob::glob;
use std::{error::Error, path::PathBuf};
use std::env;
use std::process;
use std::process::Command;

const LOCAL_REPO_PATH: &str = "/var/lib/solbuild/local";

fn copy_packages(current_dir: PathBuf) -> Result<(), Box<dyn Error>> {
    let mut packages = Vec::new();

    let search = format!("{}/*.eopkg", current_dir.to_str().unwrap());
    for entry in glob(&search).unwrap() {
        match entry {
            Ok(path) => {
                packages.push(path);
            },
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

    let current_dir = match env::current_dir() {
        Ok(dir) => dir,
        Err(e) => {
            eprintln!("Error getting current working directory: {}", e);
            process::exit(1);
        }
    };

    let command = args[1].as_ref();
    match command {
        "copy" => {
            sudo::escalate_if_needed()?;

            if let Err(e) = copy_packages(current_dir) {
                eprintln!("Error copying packages to local repo: {}", e);
                process::exit(1);
            }

            if let Err(e) = index_repo() {
                eprintln!("Error indexing local repo: {}", e);
                process::exit(1);
            }

            Ok(())
        },
        _ => {
            eprintln!("Unknown subcommand: {}", command);
            process::exit(1);
        }
    }
}
