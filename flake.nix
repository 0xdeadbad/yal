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
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
    yal = (import ./default.nix { inherit pkgs; });
  in
  {
    packages.${system}.default = yal;

    devShells.${system}.default = pkgs.mkShell {
      packages = [
        yal
      ];
      buildInputs = with pkgs; [
        go
      ];
    };
  };
}
