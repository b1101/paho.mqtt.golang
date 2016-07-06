/*
 * Copyright (c) 2014 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"errors"
	"strings"
)

//InvalidQos is the error returned when an packet is to be sent
//with an invalid Qos value
var ErrInvalidQos = errors.New("Invalid QoS")

//InvalidTopicEmptyString is the error returned when a topic string
//is passed in that is 0 length
var ErrInvalidTopicEmptyString = errors.New("Invalid Topic; empty string")

//InvalidTopicMultilevel is the error returned when a topic string
//is passed in that has the multi level wildcard in any position but
//the last
var ErrInvalidTopicMultilevel = errors.New("Invalid Topic; multi-level wildcard must be last level")

// Topic Names and Topic Filters
// The MQTT v3.1.1 spec clarifies a number of ambiguities with regard
// to the validity of Topic strings.
// - A Topic must be between 1 and 65535 bytes.
// - A Topic is case sensitive.
// - A Topic may contain whitespace.
// - A Topic containing a leading forward slash is different than a Topic without.
// - A Topic may be "/" (two levels, both empty string).
// - A Topic must be UTF-8 encoded.
// - A Topic may contain any number of levels.
// - A Topic may contain an empty level (two forward slashes in a row).
// - A TopicName may not contain a wildcard.
// - A TopicFilter may only have a # (multi-level) wildcard as the last level.
// - A TopicFilter may contain any number of + (single-level) wildcards.
// - A TopicFilter with a # will match the absense of a level
//     Example:  a subscription to "foo/#" will match messages published to "foo".

type Topic []string

func ParseTopic(rawtopic string) (t Topic, err error) {
	if err = validateTopic(rawtopic); err != nil {
		return
	}
	return strings.Split(rawtopic, "/"), nil
}

func NewTopic(levels ...string) (t Topic) {
	return levels
}

func (t Topic) String() string {
	return strings.Join(t, "/")
}

func (t Topic) IsValid() bool {
	err := validateTopic(t.String())
	return err == nil
}

func (t Topic) WildcardAt(index int) Topic {
	t[index] = "#"
	return t[:index]
}

func validateSubscribeMap(subs map[string]byte) ([]string, []byte, error) {
	var topics []string
	var qoss []byte
	for topic, qos := range subs {
		if err := validateTopicAndQos(topic, qos); err != nil {
			return nil, nil, err
		}
		topics = append(topics, topic)
		qoss = append(qoss, qos)
	}

	return topics, qoss, nil
}

func validateTopic(topic string) (err error) {
	if len(topic) == 0 {
		return ErrInvalidTopicEmptyString
	}

	levels := strings.Split(topic, "/")
	for i, level := range levels {
		if level == "#" && i != len(levels)-1 {
			return ErrInvalidTopicMultilevel
		}
	}
	return
}

func validateQos(qos byte) (err error) {
	if qos < 0 || qos > 2 {
		return ErrInvalidQos
	}
	return
}

func validateTopicAndQos(topic string, qos byte) (err error) {
	if err = validateTopic(topic); err != nil {
		return
	}
	return validateQos(qos)
}
