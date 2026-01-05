class ControlTower:
    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            print('Initializing control tower!')
        return cls._instance
    
    # def manageFlight(self, flight)


tower1 = ControlTower()
tower2 = ControlTower()

print(tower1 is tower2)