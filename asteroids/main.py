import pygame
from constants import *
from player import Player
from asteroid import Asteroid
from asteroidfield import AsteroidField
from shot import Shot

def game():
    #pygame.init()
    score = 0

    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    clock = pygame.time.Clock()


    asteroids = pygame.sprite.Group()
    updatable = pygame.sprite.Group()
    drawable = pygame.sprite.Group()
    shots = pygame.sprite.Group()

    Player.containers = (updatable, drawable)
    Asteroid.containers = (asteroids, updatable, drawable)
    AsteroidField.containers = (updatable)
    Shot.containers = (shots, updatable, drawable)

    player = Player(SCREEN_WIDTH/2, SCREEN_HEIGHT/2)
    asteroidField = AsteroidField()

    dt = 0

    while True:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                exit()

        screen.fill((0,0,0))

        updatable.update(dt)

        for object in drawable:
            object.draw(screen)
            if object == player:
                player.shot_timer -= dt

        for asteroid in asteroids:
            if player.check_collision(asteroid):
                print("Game over!")
                print(f"Your score was {score}.")
                return

        for asteroid in asteroids:
            for shot in shots:
                if shot.check_collision(asteroid):
                    if asteroid.radius <= ASTEROID_MIN_RADIUS:
                        score += 1
                    shot.kill()
                    asteroid.split()

        dt = clock.tick(60) / 1000
        pygame.display.flip()

def main():
    lives = 3
    restart = 'N'
    while True:
        pygame.init()
        game()

        restart = False
        while not restart:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    exit()
                if event.type == pygame.KEYDOWN and event.key == pygame.K_n:
                    exit()
                if event.type == pygame.KEYDOWN and event.key == pygame.K_y:
                    restart = True



if __name__ == "__main__":
    main()