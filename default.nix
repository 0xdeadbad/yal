{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation {
  name = "yal";
  src = ./.;

  # unpackPhase = ''
  #   for srcFile in $src; do
  #       cp -r $srcFile $(stripHash $srcFile)
  #   done
  # '';

  nativeBuildInputs = with pkgs; [
    go
  ];

  configurePhase = ''
    export TMPGODIR="$TMPDIR/$(stripHash $srcFile)"
    export TMPCACHEDIR="$TMPGODIR/cache"
    export TMPOUTDIR="$TMPGODIR/out"
    mkdir -p "$TMPOUTDIR"
    mkdir -p "$TMPCACHEDIR"
  '';

  buildPhase = ''
    CGO=0 GOCACHE="$TMPCACHEDIR" go build -ldflags "-s -w" -o "$TMPOUTDIR/$name" cmd/compiler/main.go
  '';

  installPhase = ''
    mkdir -p $out/bin
    cp $TMPOUTDIR/$name $out/bin
  '';
}
