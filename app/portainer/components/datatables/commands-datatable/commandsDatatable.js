angular.module('portainer.app').component('commandsDatatable', {
  templateUrl: 'app/portainer/components/datatables/commands-datatable/commandsDatatable.html',
  controller: 'GenericDatatableController',
  bindings: {
    title: '@',
    titleIcon: '@',
    dataset: '<',
    tableKey: '@',
    orderBy: '@',
    reverseOrder: '<',
    showTextFilter: '<',
    commandManagement: '<',
    accessManagement: '<',
    removeAction: '<'
  }
});
