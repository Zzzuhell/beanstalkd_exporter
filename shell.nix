{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
}: let
  goEnv = mkGoEnv {pwd = ./.;};
in
  pkgs.mkShell {
    packages = [
      goEnv
      gomod2nix
      pkgs.gotools
      pkgs.go-tools
      pkgs.golangci-lint
      pkgs.delve
    ];
    hardeningDisable = ["fortify"]; # https://wiki.nixos.org/wiki/Go#Using_cgo_on_NixOS
  }
