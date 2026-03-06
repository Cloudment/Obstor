/*
 * MinIO Cloud Storage (C) 2016, 2018 MinIO, Inc.
 * PGG Obstor, (C) 2021-2026 PGG, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import "core-js/stable"
import "bootstrap/dist/css/bootstrap.min.css"
import "./less/main.less"
import "@fortawesome/fontawesome-free/css/all.css"
import "material-design-iconic-font/dist/css/material-design-iconic-font.min.css"

import React from "react"
import { createRoot } from "react-dom/client"
import { unstable_HistoryRouter as HistoryRouter } from "react-router-dom"
import { Provider } from "react-redux"

import history from "./js/history"
import configureStore from "./js/store/configure-store"
import hideLoader from "./js/loader"
import App from "./js/App"

const store = configureStore()

const root = createRoot(document.getElementById("root"))
root.render(
  <Provider store={store}>
    <HistoryRouter history={history} basename="/obstor">
      <App />
    </HistoryRouter>
  </Provider>
)

hideLoader()
