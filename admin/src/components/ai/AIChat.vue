<template>
  <div class="ai-chat-container" :class="{ open: isOpen }">
    <div class="chat-toggle" @click="toggleChat">
      <el-icon><ChatDotRound /></el-icon>
    </div>
    <div class="chat-window" v-if="isOpen">
      <div class="chat-header">
        <span>AI Assistant</span>
        <el-icon class="close-btn" @click="toggleChat"><Close /></el-icon>
      </div>
      <div class="chat-messages">
        <div v-for="(msg, index) in messages" :key="index" class="message" :class="msg.role">
          {{ msg.content }}
        </div>
      </div>
      <div class="chat-input">
        <el-input v-model="input" placeholder="Ask AI..." @keyup.enter="sendMessage">
          <template #append>
            <el-button @click="sendMessage"><el-icon><Position /></el-icon></el-button>
          </template>
        </el-input>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isOpen = ref(false)
const input = ref('')
const messages = ref([
  { role: 'assistant', content: 'Hello! How can I help you today?' }
])

const toggleChat = () => {
  isOpen.value = !isOpen.value
}

const sendMessage = () => {
  if (!input.value) return
  messages.value.push({ role: 'user', content: input.value })
  
  // Mock AI response
  setTimeout(() => {
    messages.value.push({ role: 'assistant', content: 'I am a mock AI response. Backend integration coming soon!' })
  }, 1000)
  
  input.value = ''
}
</script>

<style scoped>
.ai-chat-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 1000;
}
.chat-toggle {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: #409EFF;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  font-size: 24px;
}
.chat-window {
  position: absolute;
  bottom: 60px;
  right: 0;
  width: 300px;
  height: 400px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  border: 1px solid #e6e6e6;
}
.chat-header {
  padding: 10px 15px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  background-color: #f5f7fa;
  border-radius: 8px 8px 0 0;
}
.close-btn {
  cursor: pointer;
}
.chat-messages {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
}
.message {
  margin-bottom: 10px;
  padding: 8px 12px;
  border-radius: 4px;
  max-width: 80%;
}
.message.user {
  background-color: #409EFF;
  color: white;
  align-self: flex-end;
  margin-left: auto;
}
.message.assistant {
  background-color: #f0f2f5;
  color: #333;
}
.chat-input {
  padding: 10px;
  border-top: 1px solid #eee;
}
</style>
