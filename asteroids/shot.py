import pygame
from circleshape import *
from constants import *

class Shot(CircleShape):
    def __init__(self, x, y, radius):
        super().__init__(x, y, radius)
        self.rotation = 0

    def draw(self, screen):
        pygame.draw.circle(surface=screen, color=(255,255,255), center=self.position, radius=SHOT_RADIUS, width=2)

    def update(self, dt):
        self.position += self.velocity * dt