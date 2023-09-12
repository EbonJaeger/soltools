use std::{
    fs::{read_dir, remove_file},
    process::Command,
};

use globset::{Glob, GlobSetBuilder};
use thiserror::Error;

const LOCAL_REPO: &str = "/var/lib/solbuild/local";

/// Globs files to remove from the local Solbuild repository.
///
/// If nothing is passed to `remove`, we will attempt to remove
/// all of the files in the repository.
///
/// Anything passed to `keep` will be kept, overriding any other factor.
///
/// If `dry_run` is set, then this function will only print the names of the
/// files that would normally be removed, without actually removing them.
///
/// The repository will be indexed after removals have been performed if
/// `index` is set.
pub fn handle(
    dry_run: bool,
    index: bool,
    remove: Option<Vec<String>>,
    keep: Option<Vec<&str>>,
) -> Result<(), Error> {
    // Build a glob set for files to remove
    let builder = match remove {
        Some(removals) => {
            let mut builder = GlobSetBuilder::new();
            for removal in removals {
                builder.add(Glob::new(&format!("{}/{}*.eopkg", LOCAL_REPO, removal))?);
            }
            builder
        }
        None => {
            let mut builder = GlobSetBuilder::new();
            builder.add(Glob::new(&format!("{}/*.eopkg", LOCAL_REPO))?);
            builder
        }
    };

    // Build a new glob set for files that should be kept
    let keeper_builder = if let Some(keepers) = keep {
        let mut keeper_builder = GlobSetBuilder::new();
        for keeper in keepers {
            keeper_builder.add(Glob::new(&format!("{}/{}*.eopkg", LOCAL_REPO, keeper))?);
        }
        Some(keeper_builder)
    } else {
        None
    };

    // Build the glob sets
    let set = builder.build()?;
    let keeper_set = if let Some(k) = keeper_builder {
        Some(k.build()?)
    } else {
        None
    };

    // Iterate over all files in the repository
    for entry in read_dir(LOCAL_REPO)? {
        let entry = entry?;
        if set.is_match(entry.path()) {
            // Check if this file should be kept
            if let Some(k) = &keeper_set {
                if k.is_match(entry.path()) {
                    continue;
                }
            }

            // Only print the file name if it's a dry run
            if dry_run {
                println!("{}", entry.path().display());
                continue;
            }

            // Remove the file
            remove_file(entry.path())?;
        }
    }

    // Index the repository
    if index && !dry_run {
        Command::new("eopkg")
            .args(["index", "--skip-signing", LOCAL_REPO])
            .current_dir(LOCAL_REPO)
            .output()?;
    }

    Ok(())
}

#[derive(Debug, Error)]
pub enum Error {
    #[error("Remove error")]
    Io(#[from] std::io::Error),
    #[error("Glob error")]
    Glob(#[from] globset::Error),
}
