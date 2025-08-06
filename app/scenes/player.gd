extends CharacterBody2D

const speed = 30
var curr_dir = "none"

func _ready():
	$AnimatedSprite2D.play("front_idle")

func _physics_process(delta):
	player_movement(delta)
	
func player_movement(delta):
	
	if Input.is_action_pressed("ui_right"):
		curr_dir = "right"
		play_anim(1)
		velocity.x = speed
		velocity.y = 0
	elif Input.is_action_pressed("ui_left"):
		curr_dir = "left"
		play_anim(1)
		velocity.x = -speed
		velocity.y = 0
	elif Input.is_action_pressed("ui_down"):
		play_anim(1)
		curr_dir = "down"
		velocity.y = speed
		velocity.x = 0
	elif Input.is_action_pressed("ui_up"):
		play_anim(1)
		curr_dir = "up"
		velocity.y = -speed
		velocity.x = 0
	else:
		play_anim(0)
		velocity.x = 0
		velocity.y = 0

	move_and_slide()

func play_anim(movement):
	var dir = curr_dir
	var anim = $AnimatedSprite2D
	
	if dir == "right" or dir=="left":
		anim.flip_h = !(dir == "right")
		if(movement == 1):
			anim.play("side_walk")
		elif(movement == 0):
			anim.play("side_idle")
	elif dir == "up":
		anim.flip_h = true
		if(movement == 1):
				anim.play("back_walk")
		elif(movement == 0):
				anim.play("back_idle")
	elif dir == "down":
		anim.flip_h = true
		if(movement == 1):
			anim.play("front_walk")
		elif(movement == 0):
			anim.play("front_idle")
			
		
