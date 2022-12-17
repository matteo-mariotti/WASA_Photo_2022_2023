<script>
export default {
props: ['open'],
}
</script>


<template>

  <div class="vue-modal" v-show="open">
    <div class="vue-modal-inner">
      <div class="vue-modal-content">
        <slot />
        <br>
        <button class="btn btn-dark" @click="$emit('close')"> Close </button>

      </div>
    </div>
  </div>

</template>

<style scoped>

  ::before,
  ::after{
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  .vue-modal{
    position: fixed;
    top: 0;
    width: 100%;
    height: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    background-color: rgba(92,92,92,0.45);
    z-index: 2;
  }

  .vue-modal-inner{
    max-width: 600px;
    margin: auto;
    margin-top: 20%;

  }

  .vue-modal-content{
    position: relative;
    background-color: white;
    border: 1px solid black;
    padding: 1rem;
    background-clip: padding-box;
    border-radius: 0.3rem;
    font-size: 1rem;
  }

</style>