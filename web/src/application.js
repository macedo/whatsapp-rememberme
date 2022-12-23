import { Application } from "@hotwired/stimulus";

import NotificationController from "./controllers/notification_controller";

window.Stimulus = Application.start();
Stimulus.register("notification", NotificationController);
