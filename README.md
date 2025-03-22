# digital-dungeon-master

Digital Dungeon Master - DDM - is an Ai driven role playing game

## Todos:

- Clear input form after send
- add a "command overwrite", instead of frontend {} have a separate input fiel for instructions
- Translate oracle tables into sql tables
- Maybe use two chats with LLM. One with history that is the GM, and one for preparing questions and instructions (i.e. correctly phrasing dice rolls and oracle table lookups)
- Get current token use / total tokens
- How to handle dice throws?
- How does the user - LLM communication work in regards to game mechanics? Where and when are dice calculated, how do LLM and user learn the results?
- Aside from the plain text LLM response, how to read out specific location and enemey data, that needs to be displayed in the game?
- How to handle health, mana, stamina and / or other stats?
- Create llamafiles and embed or at least include and automatically start
- How is character data handled?
- Delete last ai response and regenerate
- HTMX scroll to last chat message
- Create game art
- Incorporate game art
- Model Ideas: https://www.reddit.com/r/LocalLLaMA/comments/1ge19ps/llm_model_for_dnd/
- Model Ideas: https://www.drivethrurpg.com/en/product/494922/oracle-ai-roleplaying-guide
- Get a good comfy Ui setup and use it directly with flux (without comfyui)
- Dockerize flux

- Idea:
  <history>the chat up until here</history>
  <location>the current location</location>
  <player>the player's characters current stats</player>
  <enemy>the enemy's characters current stats</enemy>
  <rag-data>the rag data from the vector db</rag-data>
  <user-message>the user's communication with the game master</user-message>
