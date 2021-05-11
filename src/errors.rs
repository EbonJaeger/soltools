use fs_extra::error::Error as FsError;
use git2::Error as GitError;
use glob::GlobError;
use std::error::Error;
use std::fmt::{self, Display};
use std::io::Error as IoError;
use std::result::Result;

/// A common result type used by functions in this application.
pub type SolResult<T> = Result<T, SolError>;

#[derive(Debug)]
#[non_exhaustive]
pub enum SolError {
    Fs(FsError),
    Git(GitError),
    Glob(GlobError),
    Io(IoError),
    Other(&'static str),
}

impl From<FsError> for SolError {
    fn from(e: FsError) -> SolError {
        SolError::Fs(e)
    }
}

impl From<GitError> for SolError {
    fn from(e: GitError) -> SolError {
        SolError::Git(e)
    }
}

impl From<GlobError> for SolError {
    fn from(e: GlobError) -> SolError {
        SolError::Glob(e)
    }
}

impl From<IoError> for SolError {
    fn from(e: IoError) -> SolError {
        SolError::Io(e)
    }
}

impl Display for SolError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match &self {
            SolError::Fs(inner) => fmt::Display::fmt(&inner, f),
            SolError::Git(inner) => fmt::Display::fmt(&inner, f),
            SolError::Glob(inner) => fmt::Display::fmt(&inner, f),
            SolError::Io(inner) => fmt::Display::fmt(&inner, f),
            SolError::Other(msg) => f.write_str(msg),
        }
    }
}

impl Error for SolError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match self {
            SolError::Fs(inner) => Some(inner),
            SolError::Git(inner) => Some(inner),
            SolError::Glob(inner) => Some(inner),
            SolError::Io(inner) => Some(inner),
            _ => None,
        }
    }
}
