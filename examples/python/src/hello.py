from markdown import markdown

def hello():
    return markdown("Hello *Earthly*")

print(hello())