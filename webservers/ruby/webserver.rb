require 'socket'
require 'digest'
require 'base64'
require 'uri'

server = TCPServer.new(8008)

puts "Server started on port 8008..."

loop do
  client = server.accept

  # Read the request
  request_line = client.gets
  next unless request_line

  method, path, _ = request_line.split
  headers = {}
  while (line = client.gets) && line != "\r\n"
    key, value = line.split(': ', 2)
    headers[key] = value.strip
  end
  body = client.read(headers['Content-Length'].to_i) if headers['Content-Length']

  # Parse POST data
  params = {}
  if body
    body.split('&').each do |pair|
      key, value = pair.split('=')
      params[key] = URI.decode_www_form_component(value)
    end
  end

  # Handle routes and methods
  response = case [method, path]
  when ['POST', '/sha256']
    input = params['input']
    if input.nil? || input.empty?
      "HTTP/1.1 400 Bad Request\r\n\r\nInput is required"
    else
      "HTTP/1.1 200 OK\r\n\r\n#{Digest::SHA256.hexdigest(input)}"
    end
  when ['POST', '/base64']
    input = params['input']
    if input.nil? || input.empty?
      "HTTP/1.1 400 Bad Request\r\n\r\nInput is required"
    else
      "HTTP/1.1 200 OK\r\n\r\n#{Base64.strict_encode64(input)}"
    end
  when ['POST', '/urlencode']
    input = params['input']
    if input.nil? || input.empty?
      "HTTP/1.1 400 Bad Request\r\n\r\nInput is required"
    else
      "HTTP/1.1 200 OK\r\n\r\n#{URI.encode_www_form_component(input)}"
    end
  when ['GET', '/']
    "HTTP/1.1 201 Created\r\n\r\nruby"
  else
    "HTTP/1.1 404 Not Found\r\n\r\nNot Found"
  end

  # Send the response
  client.print response
  client.close
end
