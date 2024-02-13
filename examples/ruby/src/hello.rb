require 'colorize'

def hello
  'Hello'.colorize(:blue) + ' ' + 'Earthly'.colorize(:green)
end

puts hello
