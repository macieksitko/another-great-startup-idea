from fastapi import FastAPI, HTTPException, Body
from fastembed import TextEmbedding
from typing import List
from pydantic import BaseModel

app = FastAPI()

# Initialize the embedding model
embedding_model = TextEmbedding(model_name="BAAI/bge-small-en-v1.5")

class TextRequest(BaseModel):
    texts: List[str]

@app.get("/")
async def root():
    return {"message": "Word Embeddings API"}

@app.post("/embed")
async def get_embeddings(request_data: TextRequest):
    try:
        embeddings = list(embedding_model.embed(request_data.texts))
        return {
            "embeddings": [emb.tolist() for emb in embeddings],
            "dimensions": len(embeddings[0]) if embeddings else 0
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
