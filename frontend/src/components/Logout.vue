<template>
  <div class="logout">
    <form @submit.prevent="logout">
      <button type="logout">Logout</button>
    </form>
  </div>
</template>

<script>
export default {
    name: 'Logout',
    methods: {
        logout() {
            this.$http.post('/api/user/logout')
                .then(() => {
                    this.$emit('logout');
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
    },
};
</script>
