<template>
  <div class="home">
    <h1>Welcome {{ username }}</h1>
  </div>
</template>

<script>
export default {
    name: 'home',
    data() {
        return {
            id: '',
            username: '',
            email: '',
        };
    },
    props: { authenticated: Boolean },
    created() {
        this.$http.getProtected('/api/followings')
            .then((response) => {
                console.log(response.data);
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>
