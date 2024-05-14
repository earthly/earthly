defmodule Earthly.MixProject do
  use Mix.Project

  def project do
    [
      app: :earthly,
      version: "0.1.0",
      elixir: "~> 1.16",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      applications: [:timex],
      extra_applications: [:logger],
      mod: {Earthly.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:timex, "~> 3.0"}
    ]
  end
end
