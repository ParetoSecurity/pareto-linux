{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-B4qq2nX6bIwY8rWqcVcZSJ36kjxvvuajVOZzV1volJE";
  subPackages = ["."];
  postInstall = ''
    mv $out/bin/pareto-core $out/bin/paretosecurity
  '';
}
