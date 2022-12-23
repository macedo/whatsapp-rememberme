import { Controller } from "@hotwired/stimulus";

export default class NotificationController extends Controller {
  connect() {
    setTimeout(() => {
      this.element.remove();
    }, 3000)
  };

  close() {
    this.element.remove();
  }
};
