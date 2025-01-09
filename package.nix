{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-UPAu4eDu0WY/ktMOiqDakfcwXRj+xvOXZ7ybXpp0uAo=";
  subPackages = ["."];
  postInstall = ''
    mv $out/bin/pareto-linux $out/bin/paretosecurity
  '';
}
