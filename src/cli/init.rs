use std::{
    env,
    fs::{create_dir, write},
    path::Path,
    process::Command,
};

use thiserror::Error;

pub fn handle(name: String, url: String) -> Result<(), Error> {
    let current = env::current_dir()?;

    // Create the directory
    let package_dir = Path::new(&current).join(&name);
    create_dir(&package_dir)?;

    // Write the Makefile
    let makefile_path = Path::new(&name).join("Makefile");
    write(makefile_path, b"include ../Makefile.common\n")?;

    // Run yauto.py
    let yauto_path = Path::new(&current).join("common/Scripts/yauto.py");
    Command::new(&yauto_path)
        .arg(&url)
        .current_dir(package_dir)
        .output()?;

    Ok(())
}

#[derive(Debug, Error)]
pub enum Error {
    #[error("Write error")]
    Io(#[from] std::io::Error),
}
