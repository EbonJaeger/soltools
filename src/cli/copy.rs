use std::{fs::copy, path::Path, process::Command};

use glob::glob;
use thiserror::Error;

const LOCAL_REPO: &str = "/var/lib/solbuild/local";

/// Globs any eopkg files and copies them to the local Solbuild
/// repository directory. Once this is done, it optionally
/// indexes the local repository.
pub fn handle(index: bool) -> Result<(), Error> {
    let entries = glob("./*.eopkg")?;

    for entry in entries {
        match entry {
            Ok(path) => {
                let final_path = Path::new(LOCAL_REPO).join(&path);
                copy(path, final_path)?;
            }
            Err(e) => return Err(Error::Glob(e)),
        }
    }

    if index {
        Command::new("eopkg")
            .args(["index", "--skip-signing", LOCAL_REPO])
            .current_dir(LOCAL_REPO)
            .output()?;
    }

    Ok(())
}

#[derive(Debug, Error)]
pub enum Error {
    #[error("Write error")]
    Io(#[from] std::io::Error),
    #[error("Glob error")]
    Glob(#[from] glob::GlobError),
    #[error("Glob path error")]
    Path(#[from] glob::PatternError),
}
