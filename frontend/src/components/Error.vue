<template>
  <div class="error" v-on:click="hide">
    <div class="pointer"/>
    <div class="body">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
    name: 'Error',
    methods: {
        hide() {
            this.$el.style.visibility = "hidden";
        },
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.error {
  display: inline-block;
}

.pointer {
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-bottom: 10px solid $error;
  position: relative;
  margin: auto;
}

.body {
  box-shadow: 0 0 4px 0 $base3;
  background-color: $error;
  color: $text1;
  padding: $margin0 $margin1;
  border-radius: $margin0;

  a {
    font-weight: bold;
    color: $text1;
  }
}
</style>
