use std::{
    env,
    fs::{create_dir, write},
    path::Path,
    process::Command,
};

use thiserror::Error;

const MAINTAINER_CONTENTS: &str = r#"This file is used to indicate primary maintainership for this package. A package may list more than one maintainer to avoid bus factor issues. People on this list may be considered “subject-matter experts”. Please note that Solus staff may need to perform necessary rebuilds, upgrades, or security fixes as part of the normal maintenance of the Solus package repository. If you believe this package requires an update, follow documentation from https://help.getsol.us/docs/packaging/procedures/request-a-package-update. In the event that this package becomes insufficiently maintained, the Solus staff reserves the right to request a new maintainer, or deprecate and remove this package from the repository entirely.

- Evan Maddock
  - Matrix: @ebonjaeger:matrix.org
  - Email: maddock.evan@vivaldi.net
"#;

pub fn handle(name: String, url: String, maintain: bool) -> Result<(), Error> {
    let current = env::current_dir()?;

    // Create the directory
    let package_dir = Path::new(&current).join(&name);
    create_dir(&package_dir)?;

    // Write the Makefile
    let makefile_path = Path::new(&name).join("Makefile");
    write(makefile_path, b"include ../Makefile.common\n")?;

    // Optionally write the maintainer file
    if maintain {
        let maintainer_file_path = Path::new(&name).join("MAINTAINERS.md");
        write(maintainer_file_path, MAINTAINER_CONTENTS)?;
    }

    // Run yauto.py
    let yauto_path = Path::new(&current).join("common/Scripts/yauto.py");
    Command::new(yauto_path)
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
