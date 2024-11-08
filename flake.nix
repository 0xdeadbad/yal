{
  description = "YAL (Yet Another Language) is a simple scripting language written in Go";

  inputs =
    {
      nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    };

    outputs = {
      self,
      nixpkgs
    }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system}.default = (import ./default.nix { inherit pkgs; });

      devShells.${system}.default = pkgs.mkShell {
        buildInputs = [];
      };
    };
}
