<template>
  <v-dialog :value="show" lazy persistent max-width="500" class="dialog-album-edit" color="application"
            @keydown.esc="close">
    <v-form ref="form" lazy-validation
            dense class="form-album-edit" accept-charset="UTF-8"
            @submit.prevent="confirm">
      <v-card raised elevation="24">
        <v-card-title primary-title class="pb-0">
          <v-layout row wrap>
            <v-flex xs12>
              <h3 class="headline mx-2 mb-0">
                <translate :translate-params="{name: model.modelName()}">Edit %{name}</translate>
              </h3>
            </v-flex>
          </v-layout>
        </v-card-title>

        <v-card-text>
          <v-container fluid class="pa-0">
            <v-layout row wrap>
              <v-flex v-if="album.Type !== 'month'" xs12 pa-2>
                <v-text-field v-model="model.Title"
                              hide-details autofocus
                              :rules="[titleRule]"
                              :label="$gettext('Name')"
                              :disabled="disabled"
                              color="secondary-dark"
                              class="input-title"
                              @keyup.enter.native="confirm"
                ></v-text-field>
              </v-flex>
              <v-flex xs12 pa-2>
                <v-text-field v-model="model.Location"
                              hide-details
                              :label="$gettext('Location')"
                              :disabled="disabled"
                              color="secondary-dark"
                              class="input-location"
                ></v-text-field>
              </v-flex>
              <v-flex xs12 pa-2>
                <v-textarea :key="growDesc" v-model="model.Description"
                            auto-grow
                            hide-details
                            browser-autocomplete="off"
                            :label="$gettext('Description')"
                            :rows="1"
                            :disabled="disabled"
                            class="input-description"
                            color="secondary-dark">
                </v-textarea>
              </v-flex>
              <v-flex xs12 md6 pa-2>
                <v-combobox v-model="model.Category" hide-details
                            :search-input.sync="model.Category"
                            :items="categories"
                            :disabled="disabled"
                            :label="$gettext('Category')"
                            :allow-overflow="false"
                            return-masked-value
                            color="secondary-dark"
                            class="input-category"
                >
                </v-combobox>
              </v-flex>
              <v-flex xs12 md6 pa-2>
                <v-select
                    v-model="model.Order"
                    :label="$gettext('Sort Order')"
                    hide-details
                    :items="sorting"
                    :disabled="disabled"
                    item-value="value"
                    item-text="text"
                    color="secondary-dark">
                </v-select>
              </v-flex>
            </v-layout>
          </v-container>
        </v-card-text>
        <v-card-actions class="pt-0">
          <v-layout row wrap class="pa-2">
            <v-flex xs12 text-xs-right>
              <v-btn depressed color="secondary-light"
                     class="action-cancel"
                     @click.stop="close">
                <translate>Cancel</translate>
              </v-btn>
              <v-btn depressed dark color="primary-button"
                     class="action-confirm"
                     :disabled="disabled"
                     @click.stop="confirm">
                <translate>Save</translate>
              </v-btn>
            </v-flex>
          </v-layout>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
<script>
import Album from "model/album";

export default {
  name: 'PAlbumEditDialog',
  props: {
    show: Boolean,
    album: {
      type: Object,
      default: () => {},
    },
  },
  data() {
    return {
      disabled: !this.$config.allow("albums", "manage"),
      model: new Album(),
      growDesc: false,
      loading: false,
      sorting: [
        {value: 'added', text: this.$gettext('Recently added')},
        {value: 'edited', text: this.$gettext('Recently edited')},
        {value: 'newest', text: this.$gettext('Newest first')},
        {value: 'oldest', text: this.$gettext('Oldest first')},
        {value: 'name', text: this.$gettext('Sort by file name')},
        {value: 'similar', text: this.$gettext('Group by similarity')},
        {value: 'relevance', text: this.$gettext('Most relevant')},
      ],
      categories: this.$config.albumCategories(),
      titleRule: v => v.length <= this.$config.get('clip') || this.$gettext("Name too long"),
    };
  },
  watch: {
    show: function (show) {
      if (show) {
        this.model = this.album.clone();
      }
    }
  },
  methods: {
    expand() {
      this.growDesc = !this.growDesc;
    },
    close() {
      this.$emit('close');
    },
    confirm() {
      if (this.disabled) {
        this.close();
        return;
      }

      this.model.update().then((m) => {
        this.categories = this.$config.albumCategories();
        this.$emit('close');
      });
    },
  },
};
</script>
