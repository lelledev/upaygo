{ pkgs, ... }: {

  # Which nixpkgs channel to use.
  channel = "stable-23.11"; # or "unstable"

  # Use https://search.nixos.org/packages to find packages
  packages = [
    pkgs.go
    pkgs.docker
  ];

  # Sets environment variables in the workspace
  env = { };

  # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
  idx.extensions = [
    "open-vsx.go"
  ];

  idx.workspace.onCreate = {
    go-tidy = "go mod tidy";
    go-vendor = "go mod vendor";
  };

  idx.workspace.onStart = {
    go-tidy = "go mod tidy";
    go-vendor = "go mod vendor";
  };

  services.docker.enable = true;

  # Enable previews and customize configuration
  idx.previews = {
    enable = false;
    previews = [ ];
  };
}
