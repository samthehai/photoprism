<template>
  <v-form ref="form" lazy-validation
          dense autocomplete="off" class="p-photo-toolbar" accept-charset="UTF-8"
          @submit.prevent="updateQuery()">
    <v-toolbar flat :dense="$vuetify.breakpoint.smAndDown" class="page-toolbar" color="secondary">
      <v-text-field :value="filter.q"
                    class="input-search background-inherit elevation-0"
                    solo hide-details clearable overflow single-line
                    validate-on-blur
                    autocorrect="off"
                    autocapitalize="none"
                    browser-autocomplete="off"
                    :label="$gettext('Search')"
                    prepend-inner-icon="search"
                    color="secondary-dark"
                    @change="(v) => {updateFilter({'q': v})}"
                    @keyup.enter.native="(e) => updateQuery({'q': e.target.value})"
                    @click:clear="() => {updateQuery({'q': ''})}"
      ></v-text-field>

      <v-btn icon class="hidden-xs-only action-reload" :title="$gettext('Reload')" @click.stop="refresh()">
        <v-icon>refresh</v-icon>
      </v-btn>

      <v-btn v-if="settings.view === 'cards'" icon :title="$gettext('Toggle View')" @click.stop="setView('list')">
        <v-icon>view_list</v-icon>
      </v-btn>
      <v-btn v-else-if="settings.view === 'list'" icon :title="$gettext('Toggle View')" @click.stop="setView('mosaic')">
        <v-icon>view_comfy</v-icon>
      </v-btn>
      <v-btn v-else icon :title="$gettext('Toggle View')" @click.stop="setView('cards')">
        <v-icon>view_column</v-icon>
      </v-btn>

      <v-btn v-if="!$config.values.readonly && $config.feature('upload')" icon class="hidden-sm-and-down action-upload"
             :title="$gettext('Upload')" @click.stop="showUpload()">
        <v-icon>cloud_upload</v-icon>
      </v-btn>

      <v-btn icon class="p-expand-search" :title="$gettext('Expand Search')"
             @click.stop="searchExpanded = !searchExpanded">
        <v-icon>{{ searchExpanded ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}</v-icon>
      </v-btn>
    </v-toolbar>

    <v-card v-show="searchExpanded"
            class="pt-1"
            flat
            color="secondary-light">
      <v-card-text>
        <v-layout row wrap>
          <v-flex xs12 sm6 md3 pa-2 class="p-countries-select">
            <v-select :value="filter.country"
                      :label="$gettext('Country')"
                      flat solo hide-details
                      color="secondary-dark"
                      item-value="ID"
                      item-text="Name"
                      :items="countryOptions"
                      class="input-countries"
                      @change="(v) => {updateQuery({'country': v})}"
            >
            </v-select>
          </v-flex>
          <v-flex xs12 sm6 md3 pa-2 class="p-camera-select">
            <v-select :value="filter.camera"
                      :label="$gettext('Camera')"
                      flat solo hide-details
                      color="secondary-dark"
                      item-value="ID"
                      item-text="Name"
                      :items="cameraOptions"
                      @change="(v) => {updateQuery({'camera': v})}">
            </v-select>
          </v-flex>
          <v-flex xs12 sm6 md3 pa-2 class="p-view-select">
            <v-select id="viewSelect"
                      :value="settings.view"
                      :label="$gettext('View')" flat solo
                      hide-details
                      color="secondary-dark"
                      :items="options.views"
                      @change="(v) => {setView(v)}">
            </v-select>
          </v-flex>
          <v-flex xs12 sm6 md3 pa-2 class="p-time-select">
            <v-select :value="filter.order"
                      :label="$gettext('Sort Order')"
                      flat solo hide-details
                      color="secondary-dark"
                      :items="options.sorting"
                      @change="(v) => {updateQuery({'order': v})}">
            </v-select>
          </v-flex>
          <v-flex xs12 sm6 md3 pa-2 class="p-year-select">
            <v-select :value="filter.year"
                      :label="$gettext('Year')"
                      flat solo hide-details
                      color="secondary-dark"
                      item-value="value"
                      item-text="text"
                      :items="yearOptions()"
                      @change="(v) => {updateQuery({'year': v})}">
            </v-select>
          </v-flex>
          <v-flex xs12 sm6 md3 pa-2 class="p-month-select">
            <v-select :value="filter.month"
                      :label="$gettext('Month')"
                      flat solo hide-details
                      color="secondary-dark"
                      item-value="value"
                      item-text="text"
                      :items="monthOptions()"
                      @change="(v) => {updateQuery({'month': v})}">
            </v-select>
          </v-flex>
          <!-- v-flex xs12 sm6 md3 pa-2 class="p-lens-select">
              <v-select @change="dropdownChange"
                        :label="labels.lens"
                        flat solo hide-details
                        color="secondary-dark"
                        item-value="ID"
                        item-text="Model"
                        v-model="filter.lens"
                        :items="lensOptions">
              </v-select>
          </v-flex -->
          <v-flex xs12 sm6 md3 pa-2 class="p-color-select">
            <v-select :value="filter.color"
                      :label="$gettext('Color')"
                      flat solo hide-details
                      color="secondary-dark"
                      item-value="Slug"
                      item-text="Name"
                      :items="colorOptions()"
                      @change="(v) => {updateQuery({'color': v})}">
            </v-select>
          </v-flex>
          <v-flex xs12 sm6 md3 pa-2 class="p-category-select">
            <v-select :value="filter.label"
                      :label="$gettext('Category')"
                      flat solo hide-details
                      color="secondary-dark"
                      item-value="Slug"
                      item-text="Name"
                      :items="categoryOptions"
                      @change="(v) => {updateQuery({'label': v})}">
            </v-select>
          </v-flex>
        </v-layout>
      </v-card-text>
    </v-card>
  </v-form>
</template>
<script>
import Event from "pubsub-js";
import * as options from "options/options";

export default {
  name: 'PPhotoToolbar',
  props: {
    filter: {
      type: Object,
      default: () => {},
    },
    updateFilter: {
      type: Function,
      default: () => {},
    },
    updateQuery: {
      type: Function,
      default: () => {},
    },
    settings: {
      type: Object,
      default: () => {},
    },
    refresh: {
      type: Function,
      default: () => {},
    },
  },
  data() {
    return {
      experimental: this.$config.get("experimental"),
      isFullScreen: !!document.fullscreenElement,
      config: this.$config.values,
      searchExpanded: false,
      all: {
        countries: [{ID: "", Name: this.$gettext("All Countries")}],
        cameras: [{ID: 0, Name: this.$gettext("All Cameras")}],
        lenses: [{ID: 0, Name: this.$gettext("All Lenses")}],
        colors: [{Slug: "", Name: this.$gettext("All Colors")}],
        categories: [{Slug: "", Name: this.$gettext("All Categories")}],
        months: [{value: 0, text: this.$gettext("All Months")}],
        years: [{value: 0, text: this.$gettext("All Years")}],
      },
      options: {
        'views': [
          {value: 'mosaic', text: this.$gettext('Mosaic')},
          {value: 'cards', text: this.$gettext('Cards')},
          {value: 'list', text: this.$gettext('List')},
        ],
        'sorting': [
          {value: 'added', text: this.$gettext('Recently added')},
          {value: 'edited', text: this.$gettext('Recently edited')},
          {value: 'newest', text: this.$gettext('Newest first')},
          {value: 'oldest', text: this.$gettext('Oldest first')},
          {value: 'name', text: this.$gettext('Sort by file name')},
          {value: 'similar', text: this.$gettext('Group by similarity')},
          {value: 'relevance', text: this.$gettext('Most relevant')},
          {value: 'duration', text: this.$gettext('Duration')},
        ],
      },
    };
  },
  computed: {
    countryOptions() {
      return this.all.countries.concat(this.config.countries);
    },
    cameraOptions() {
      return this.all.cameras.concat(this.config.cameras);
    },
    lensOptions() {
      return this.all.lenses.concat(this.config.lenses);
    },
    categoryOptions() {
      return this.all.categories.concat(this.config.categories);
    },
  },
  methods: {
    colorOptions() {
      return this.all.colors.concat(options.Colors());
    },
    monthOptions() {
      return this.all.months.concat(options.Months());
    },
    yearOptions() {
      return this.all.years.concat(options.IndexedYears());
    },
    setView(name) {
      if (name) {
        this.refresh({'view': name});
      }
    },
    showUpload() {
      Event.publish("dialog.upload");
    }
  },
};
</script>
