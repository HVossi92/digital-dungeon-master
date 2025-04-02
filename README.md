# Digital Dungeon Master (DDM)

Digital Dungeon Master (DDM) is an AI-powered role-playing game assistant that brings your D&D campaigns to life. It acts as an intelligent, adaptive Dungeon Master that responds to player actions, narrates adventures, and manages game mechanics - all powered by large language models.

## Overview

DDM creates an immersive D&D experience by providing an AI Dungeon Master that:
- Narrates rich, dynamic storylines that adapt to player choices
- Manages combat encounters and NPC interactions
- Handles game mechanics like dice rolls and skill checks
- Maintains consistent world state and character information
- Generates atmospheric descriptions and engaging dialogue

## Features

- **Interactive Chat Interface**: Communicate naturally with the AI Dungeon Master
- **Context-Aware Responses**: The AI maintains campaign history and character states
- **Structured Game Data**: Tracks locations, player stats, enemy information, and more
- **Dice Rolling System**: Fair and transparent dice mechanics
- **Command System**: Special commands for game mechanics and meta-actions

## Important Usage Notes

When playing with DDM, keep these guidelines in mind:

- **Play in Good Faith**: While the AI can be manipulated, the experience is most enjoyable when players engage honestly with the story
- **Clear Communication**: Be specific in your actions and intentions to help the AI understand
- **Respect Game Mechanics**: Use proper commands for dice rolls and game actions
- **Stay in Character**: Immerse yourself in role-playing for the best experience
- **Address NPCs directly**: Make it very clear to the AI who you are talking to

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running
- NVIDIA GPU recommended for optimal performance
- Docker (for production deployment)
- Go for local development

## Tech Stack

- **Frontend**: 
  - HTMX for dynamic interactions
  - Bootstrap CSS for styling
  
- **Backend**:
  - Go server
  - LLM integration via Ollama
  
- **AI/ML**:
  - Large Language Models for game mastering

## Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/digital-dungeon-master.git
   cd digital-dungeon-master
   ```

2. Install dependencies:
   ```bash
   # Backend dependencies
   go mod download
   ```

3. Start Ollama and ensure it's accessible

4. Start the development server:
   ```bash
   # Start backend
  air 
   ```

5. Visit `http://localhost:3000` in your browser

## Production Deployment

1. Build the Docker image:
   ```bash
   make build
   ```
## Roadmap

- [ ] Command overwrite system
- [ ] Dual LLM chat system for GM and mechanics
- [ ] Token usage tracking
- [ ] Enhanced dice mechanics
- [ ] Structured data extraction
- [ ] Llamafile integration
- [ ] HTMX improvements
- [ ] Game art generation
- [ ] Comfy UI integration
- [ ] Full containerization