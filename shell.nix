let nixpkgs = <nixpkgs>;
    config = {};
in with import nixpkgs { inherit config; }; let
  shellWrapper = writeScript "shell-wrapper" ''
    #! ${stdenv.shell}
    set -e

    exec -a shell ${fish}/bin/fish --login --interactive "$@"
  '';

  gomarkdoc = buildGoModule {
    name = "gomarkdoc";
    src = fetchFromGitHub {
      owner = "princjef";
      repo = "gomarkdoc";
      rev = "cc78abbcb570e329eb145854481acfcb1072f307";
      hash = "sha256-KLMec5rtTTwz4c+6ZqAgREPXlWYtuB7rH1ZvQKmcA9U=";
    };
    doCheck = false;
    vendorHash = "sha256-LfovwcipO3/ovHLDSLRhHcEocbKdW399o6mJ45GavBM=";
  };
in stdenv.mkDerivation rec {
  name = "nix-shell";
  buildInputs = [
    glibcLocales bashInteractive man
    nix cacert curl utillinux coreutils
    git jq yq-go tmux findutils gnumake
    go gopls golangci-lint gomarkdoc
  ];
  shellHook = ''
    export root=$(pwd)

    export LANG="en_US.UTF-8"
    export NIX_PATH="nixpkgs=${nixpkgs}"

    if [ ! -z "$PS1" ]
    then
      export SHELL="${shellWrapper}"
      exec "$SHELL"
    fi
  '';
}
