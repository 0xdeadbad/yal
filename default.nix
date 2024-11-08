{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation {
  name = "yal";
  src = ./.;

  # unpackPhase = ''
  #   for srcFile in $src; do
  #       cp -r $srcFile $(stripHash $srcFile)
  #   done
  # '';

  buildInputs = with pkgs; [
    go
  ];

  buildPhase = ''
    mkdir -p $out/go-build-tmp
    mkdir -p $out/result
    GOCACHE=$out/go-build-tmp go build -o $out/result/yal $src/cmd/compiler/main.go
    rm -rf $out/go-build-tmp
  '';

  installPhase = ''
    mkdir -p $out/bin
    cp $out/result/yal $out/bin
    rm $out/result/yal
    rmdir $out/result
  '';
}
