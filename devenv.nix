{ pkgs, lib, config, inputs, ... }: {
  languages.go.enable = true;

  packages = [ pkgs.git pkgs.gomod2nix ];

  git-hooks.hooks = {
    govet = {
      enable = true;
      pass_filenames = false;
    };
    gotest.enable = true;
    golangci-lint = {
      enable = true;
      pass_filenames = false;
    };
  };

  outputs =
    let
      name = "go2lunch";
      version = "1.0.0";
    in
    {
      app = import ./default.nix { inherit pkgs name version; };
    };

  # See full reference at https://devenv.sh/reference/options/
}
