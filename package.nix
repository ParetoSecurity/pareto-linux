{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-hnVeKOMMLnPVxLwMUxpjGjIro3fd6PUe4RUdYr2bVYs=";
  subPackages = ["."];
  postInstall = ''
    mv $out/bin/pareto-core $out/bin/paretosecurity
  '';
}
