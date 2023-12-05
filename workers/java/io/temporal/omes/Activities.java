package io.temporal.omes;

import io.temporal.activity.ActivityInterface;
import io.temporal.activity.ActivityMethod;

@ActivityInterface
public interface Activities {
  @ActivityMethod(name = "noop")
  void noopActivity();

  @ActivityMethod(name = "echo")
  String echo(String input);

  @ActivityMethod(name = "delay")
  void delay(int milliseconds) throws InterruptedException;
}
