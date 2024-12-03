<template>
  <div class="voice-recorder w-full">
    <h1>Voice Recorder</h1>   
    <div class="flex justify-center space-x-4 w-full">
      <div>
        <button 
          @click="startRecording" 
          :disabled="isRecording"
          class="btn btn-primary"
        >
          Start Recording
        </button>
      </div>
      <div>
        <button 
          @click="pauseRecording" 
          :disabled="!isRecording || isPaused"
          class="btn btn-warning"
        >
          Pause
        </button>
      </div>
      <div>
        <button 
          @click="resumeRecording" 
          :disabled="!isPaused"
          class="btn btn-success"
        >
          Resume
        </button>
      </div>
      <div>
        <button 
          @click="stopRecording" 
          :disabled="!isRecording"
          class="btn btn-danger"
        >
          Stop
        </button>
      </div>
    </div>

    <div class="status">
      <p>Status: {{ recordingStatus }}</p>
      <p v-if="isRecording">Recording Duration: {{ formatRecordingTime(recordingDuration) }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, onBeforeUnmount } from 'vue';

  const mediaRecorder = ref<MediaRecorder | null>(null)
  const audioChunks = ref<Blob[]>([])
  const isRecording = ref(false)
  const isPaused = ref(false)
  const recordingDuration = ref(0)
  const webSocket = ref<WebSocket | null>(null)
  let recordingTimer: number | undefined

  // Computed property for recording status
  const recordingStatus = computed(() => {
    if (isRecording.value && isPaused.value) return 'Paused'
    if (isRecording.value) return 'Recording'
    return 'Stopped'
  })

  // Format time to MM:SS
  const formatRecordingTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  // Start recording
  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      
      mediaRecorder.value = new MediaRecorder(stream)
      
      mediaRecorder.value.ondataavailable = (event: BlobEvent) => {
        if (event.data.size > 0) {
          audioChunks.value.push(event.data)
          sendAudioChunk(event.data)
        }
      }

      // Establish WebSocket connection
      webSocket.value = new WebSocket('ws://localhost:8080/stream')
      
      webSocket.value.onopen = () => {
        if (mediaRecorder.value) {
          mediaRecorder.value.start(1000) // Collect data every second
          isRecording.value = true
          startTimer()
        }
      }

      webSocket.value.onerror = (error: Event) => {
        console.error('WebSocket Error:', error)
      }
    } catch (error) {
      console.error('Error starting recording:', error)
      alert('Could not access microphone')
    }
  }

  // Pause recording
  const pauseRecording = () => {
    if (mediaRecorder.value && mediaRecorder.value.state === 'recording') {
      mediaRecorder.value.pause()
      isPaused.value = true
      clearInterval(recordingTimer)
    }
  }

  // Resume recording
  const resumeRecording = () => {
    if (mediaRecorder.value && mediaRecorder.value.state === 'paused') {
      mediaRecorder.value.resume()
      isPaused.value = false
      startTimer()
    }
  }

  // Stop recording
  const stopRecording = () => {
    if (mediaRecorder.value) {
      mediaRecorder.value.stop()
      isRecording.value = false
      isPaused.value = false
      clearInterval(recordingTimer)
      
      if (webSocket.value) {
        webSocket.value.close()
      }
      
      // Reset recording duration
      recordingDuration.value = 0
    }
  }

  // Send audio chunk via WebSocket
  const sendAudioChunk = (chunk: Blob) => {
    if (webSocket.value && webSocket.value.readyState === WebSocket.OPEN) {
      webSocket.value.send(chunk)
    }
  }

  // Start recording timer
  const startTimer = () => {
    recordingTimer = setInterval(() => {
      if (!isPaused.value) {
        recordingDuration.value++
      }
    }, 1000) as unknown as number
  }

  // Cleanup on component unmount
  onBeforeUnmount(() => {
    if (mediaRecorder.value) {
      stopRecording()
    }
  })
</script>
