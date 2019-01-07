<template>
  <div class="error" v-on:click="hide">
    <div class="pointer up"/>
    <div class="pointer right"/>
    <div class="pointer left"/>
    <div class="body">
      <slot></slot>
    </div>
    <div class="pointer down"/>
  </div>
</template>

<script>
export default {
    name: 'Error',
    props: {
        direction: {
            type: String,
            default: 'up',
        },
    },
    methods: {
        hide() {
            this.$el.style.display = 'none';
        },
    },
    mounted() {
        /* eslint-disable no-param-reassign */
        this.$el.querySelectorAll('.pointer').forEach((c) => {
            c.style.visibility = 'hidden';
        });

        this.$el.querySelector(`.pointer.${this.direction}`).style.visibility = 'visible';
        /* eslint-enable no-param-reassign */
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.error {
  z-index: 100;
  display: inline-block;
}

.pointer {
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-bottom: 10px solid $error;
  position: absolute;
  margin: auto;

  &.up {
    left: 0;
    right: 0;
    top: -8px;
  }

  &.down {
    transform: rotate(180deg);
    left: 0;
    right: 0;
    bottom: -8px;
  }

  &.left {
    transform: rotate(-90deg);
    left: -8px;
    top: 0;
    bottom: 0;
  }

  &.right {
    transform: rotate(90deg);
    right: -8px;
    top: 0;
    bottom: 0;
  }
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
