{
  description = "YAL (Yet Another Language) is a simple scripting language written in Go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs
  }:
  let
    forAllSystems = nixpkgs.lib.genAttrs nixpkgs.lib.platforms.unix;

    nixpkgsFor = forAllSystems (system: import nixpkgs {
      inherit system;
      config = { };
      overlays = [ ];
    });
  in
  {
    packages = forAllSystems (system:
    let
      pkgs = nixpkgsFor."${system}";
      yal = (import ./default.nix { inherit pkgs; });
    in
    {
      default = yal;
    });

    devShells = forAllSystems (system:
    let
      pkgs = nixpkgsFor."${system}";
      yal = (import ./default.nix { inherit pkgs; });
    in
    {
      default = pkgs.mkShell {
        packages = [
          yal
        ];
      };
    });
  };
}
