extends Node

var tcp := StreamPeerTCP.new()
var connected := false
var my_id := str(randi())
var players := {} # dictionary of remote players (id -> Node2D)

func _ready():
	var err = tcp.connect_to_host("127.0.0.1", 9000)
	if err == OK:
		connected = true
		print("✅ Connected")
	else:
		print("❌ Connect error:", err)

func _process(_delta):
	if not connected:
		return

	tcp.poll()

	# Read all available bytes
	while tcp.get_available_bytes() > 0:
		var line := tcp.get_utf8_string(tcp.get_available_bytes())
		for packet in line.split("\n"):
			if packet.strip_edges() == "":
				continue
			var data = JSON.parse_string(packet)
			if typeof(data) == TYPE_DICTIONARY:
				if data.has("id") and data["id"] != my_id:
					_update_or_create_remote_player(data)

func send_player_data(id, pos: Vector2, speed):
	if not connected:
		return

	var data = {
		"id": str(id),
		"x": pos.x,
		"y": pos.y,
		"speed": speed
	}
	var json_str = JSON.stringify(data) + "\n"
	var sent = tcp.put_data(json_str.to_utf8_buffer())
	if sent != OK:
		print("Send error:", sent)

func _update_or_create_remote_player(data: Dictionary) -> void:
	var pid = str(data["id"])
	if not players.has(pid):
		# spawn a new player scene (replace with your Player.tscn or a simple sprite)
		var new_player := Node2D.new()
		var sprite := Sprite2D.new()
		sprite.texture = preload("res://icon.svg") # temp visual
		new_player.add_child(sprite)
		add_child(new_player)
		players[pid] = new_player

	# update position
	players[pid].position = Vector2(data["x"], data["y"])
