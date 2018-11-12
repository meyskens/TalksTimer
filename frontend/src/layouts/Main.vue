<style lang="scss">
@import "../scss/materialize.scss";
</style>

<template>
  <div>
    <div>
      <nav>
        <div class="nav-wrapper">
          <a href="#" class="brand-logo">TalksTimer</a>
          <ul id="nav-mobile" class="right hide-on-med-and-down">
            <li><router-link to="/">Home</router-link></li>
            <li><router-link to="/about">Test</router-link></li>
          </ul>
        </div>
      </nav>
    </div>
    <slot></slot>
  </div>
</template>

