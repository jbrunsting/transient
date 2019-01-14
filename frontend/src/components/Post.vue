<template>
  <div class="post">
    <a v-if="post.postUrl" :href="post.postUrl"><h3>{{ post.title }}</h3></a>
    <h3 v-else>{{ post.title }}</h3>
    <p>{{ post.content }}</p>
    <p>{{ new Date(post.time).toLocaleString() }}</p>
  </div>
</template>

<script>
export default {
    name: 'post',
    props: { post: Object },
};
</script>
