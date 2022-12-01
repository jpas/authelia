{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils, ... }:
    let
      inherit (nixpkgs) lib;

      eachDefaultSystem = f: utils.lib.eachDefaultSystem
        (system: f system (import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        }));
    in
    with lib;
    (eachDefaultSystem (system: pkgs: {
      packages = rec {
        inherit (pkgs) authelia;
        default = authelia;
      };

      devShells.default = pkgs.mkShell {
        buildInputs = lib.attrValues {
          inherit (pkgs) go;
        };
      };
    })) // {
      overlays.default = final: prev: {
        authelia = prev.callPackage ./package.nix { };
      };
    };
}
