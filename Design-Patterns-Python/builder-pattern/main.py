class Server:
    def __init__(self):
        self.cpu = None
        self.memory = None
        self.storage = None
        self.os = None
    
    def __str__(self):
        return (
            f"Server(cpu={self.cpu}, memory={self.memory}, "
            f"storage={self.storage}, os={self.os})"
        )

class ServerBuilder:
    def __init__(self):
        self.server = Server()

    def set_cpu(self, cpu: str):
        self.server.cpu = cpu
        return self
    
    def set_memory(self, memory: str):
        self.server.memory = memory
        return self
    
    def set_storage(self, storage: str):
        self.server.storage = storage
        return self
    
    def set_os(self, os: str):
        self.server.os = os
        return self
    
    def build(self) -> Server:
        return self.server
    

def main():
    server = (
        ServerBuilder()
        .set_cpu("8-core")
        .set_memory("32GB")
        .set_storage("1TB SSD")
        .set_os("Linux")
        .build()
    )

    print(server)

if __name__ == "__main__":
    main()