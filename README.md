# digital-dungeon-master

Digital Dungeon Master - DDM - is an Ai driven role playing game

## Todos:

- Setup database
  - Vector table for game rules
  - Settings table
  - Game chat (save every user request & ai response in plain text)
- Checkout https://docs.turso.tech/features/ai-and-embeddings#how-it-works for vector improvements
- How is character data handled?
- Setup endpoint to retrieve character and location data
- Setup endpoint to retrieve character and location images
- Setup LLM communication, teach it to be a GM, add RAG context
- Setup basic chat
- How to handle dice throws?
- How does the user - LLM communication work in regards to game mechanics? Where and when are dice calculated, how do LLM and user learn the results?
- Aside from the plain text LLM response, how to read out specific location and enemey data, that needs to be displayed in the game?
- How to handle health, mana, stamina and / or other stats?
- Create llamafiles and embed or at least include and automatically start

- Idea:
  <history>the chat up until here</history>
  <location>the current location</location>
  <player>the player's characters current stats</player>
  <enemy>the enemy's characters current stats</enemy>
  <rag-data>the rag data from the vector db</rag-data>
  <user-message>the user's communication with the game master</user-message>
