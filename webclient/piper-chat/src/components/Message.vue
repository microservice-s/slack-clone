<template>
  <div class="message">
    <li class="message-list-item">
      <!--<h5 class="message-author-name">{{ message.authorName }}</h5>-->
      <div class="message-time">
        {{ message.createdAt }}
      </div>
      <div class="message-text">{{ message.body }}</div>
    </li>
  </div>
</template>
<script>
  export default {
    name: 'Message',
    props: {
      message: Object
    }
  }
</script>
