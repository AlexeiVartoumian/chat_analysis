from dotenv import load_dotenv
from openai import OpenAI
import os
from connector import *



load_dotenv()

OPEN_API_KEY = os.environ.get("OPEN_API_KEY")
client = OpenAI(api_key=OPEN_API_KEY)


sample_query = "how would I create a database in aws"

response = client.embeddings.create(
                    input = sample_query, 
                    model = "text-embedding-3-small"
                    )

vector = response.data[0].embedding

print("this is my vector")
print(vector)

# print("\n\n\n")
# new_vector = np.array(vector)

with get_connection() as conn:
    get_embeddings(conn ,vector )