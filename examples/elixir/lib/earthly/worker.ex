defmodule Earthly.Worker do
  use GenServer
  use Timex

  def start_link(state) do
    GenServer.start_link(__MODULE__, state, name: __MODULE__)
  end

  def init(_) do
    :timer.send_interval(1000, :tick)
    {:ok, 0}
  end

  def handle_info(:tick, state) do
    time = Timex.format!(Timex.now, "%H:%M:%S", :strftime)

    case state do
      0 ->
        IO.puts("Hello Earthly! The time is: #{time}")
      1 ->
        IO.puts("Hello (again) Earthly! The time is: #{time}")
      x ->
        IO.puts("Hello Earthly! I have said this #{x} times! The time is: #{time}")
    end

    {:noreply, state + 1}
  end
end
