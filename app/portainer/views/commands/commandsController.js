angular.module('portainer.app')
.controller('CommandsController', ['$scope', '$state', '$filter',  'CommandService', 'Notifications', 'SystemService', 'CommandProvider',
function ($scope, $state, $filter, CommandService, Notifications, SystemService, CommandProvider) {
  $scope.state = {
    uploadInProgress: false,
    actionInProgress: false
  };

  $scope.formValues = {
    Name: '',
    Image: '',
    Command: '',
  };

  $scope.addCommand = function() {
    var name = $scope.formValues.Name;
    var image = $scope.formValues.Image;
    var command = $scope.formValues.Command;

    var commandId;
    $scope.state.actionInProgress = true;
    CommandService.createCommand(name, image, command)
    .then(function success(data) {
      commandId = data.Id;
      SystemService.info()
      .then(function success() {
        Notifications.success('Cmmand created', name);
        $state.reload();
      })
      .catch(function error(err) {
        Notifications.error('Failure', err, 'Unable to create command');
        CommandService.deleteCommand(commandId);
      })
      .finally(function final() {
        $scope.state.actionInProgress = false;
      });
    }, function error(err) {
      $scope.state.uploadInProgress = false;
      $scope.state.actionInProgress = false;
      Notifications.error('Failure', err, 'Unable to create command');
    }, function update(evt) {
      if (evt.upload) {
        $scope.state.uploadInProgress = evt.upload;
      }
    });
  };

  $scope.removeAction = function (selectedItems) {
    var actionCount = selectedItems.length;
    angular.forEach(selectedItems, function (command) {
      CommandService.deleteCommand(command.Id)
      .then(function success() {
        Notifications.success('Command successfully removed', command.Command);
        var index = $scope.commands.indexOf(command);
        $scope.commands.splice(index, 1);
      })
      .catch(function error(err) {
        Notifications.error('Failure', err, 'Unable to remove command');
      })
      .finally(function final() {
        --actionCount;
        if (actionCount === 0) {
          $state.reload();
        }
      });
    });
  };

  function fetchCommands() {
    CommandService.commands()
    .then(function success(data) {
      $scope.commands = data;
    })
    .catch(function error(err) {
      Notifications.error('Failure', err, 'Unable to retrieve commands');
      $scope.commands = [];
    });
  }

  fetchCommands();
}]);
